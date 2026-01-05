package main

import (
	"context"
	"log"

	"github.com/yourusername/cron-observer/backend/internal/database"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Connect to database
	log.Println("Connecting to MongoDB...")
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Get collections
	projectsCollection := db.DB.Collection(database.CollectionProjects)
	tasksCollection := db.DB.Collection(database.CollectionTasks)
	taskGroupsCollection := db.DB.Collection(database.CollectionTaskGroups)
	executionsCollection := db.DB.Collection(database.CollectionExecutions)

	// Step 1: Get all projects
	var projects []models.Project
	cursor, err := projectsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to find projects: %v", err)
	}
	if err := cursor.All(ctx, &projects); err != nil {
		log.Fatalf("Failed to decode projects: %v", err)
	}

	if len(projects) == 0 {
		log.Println("No projects found. Nothing to clean up.")
		return
	}

	// Keep the first project, delete the rest
	projectToKeep := projects[0]
	log.Printf("Keeping project: %s (ID: %s)", projectToKeep.Name, projectToKeep.ID.Hex())

	if len(projects) > 1 {
		projectIDsToDelete := make([]primitive.ObjectID, 0, len(projects)-1)
		for i := 1; i < len(projects); i++ {
			projectIDsToDelete = append(projectIDsToDelete, projects[i].ID)
		}

		// Delete other projects
		deleteResult, err := projectsCollection.DeleteMany(ctx, bson.M{
			"_id": bson.M{"$in": projectIDsToDelete},
		})
		if err != nil {
			log.Fatalf("Failed to delete projects: %v", err)
		}
		log.Printf("Deleted %d projects", deleteResult.DeletedCount)
	}

	// Step 2: Get all task groups for the kept project
	var taskGroups []models.TaskGroup
	cursor, err = taskGroupsCollection.Find(ctx, bson.M{"project_id": projectToKeep.ID})
	if err != nil {
		log.Fatalf("Failed to find task groups: %v", err)
	}
	if err := cursor.All(ctx, &taskGroups); err != nil {
		log.Fatalf("Failed to decode task groups: %v", err)
	}

	var taskGroupToKeep *models.TaskGroup
	if len(taskGroups) == 0 {
		log.Println("No task groups found for the project.")
	} else {
		// Keep the first task group, delete the rest
		taskGroupToKeep = &taskGroups[0]
		log.Printf("Keeping task group: %s (ID: %s)", taskGroupToKeep.Name, taskGroupToKeep.ID.Hex())

		if len(taskGroups) > 1 {
			taskGroupIDsToDelete := make([]primitive.ObjectID, 0, len(taskGroups)-1)
			for i := 1; i < len(taskGroups); i++ {
				taskGroupIDsToDelete = append(taskGroupIDsToDelete, taskGroups[i].ID)
			}

			deleteResult, err := taskGroupsCollection.DeleteMany(ctx, bson.M{
				"_id": bson.M{"$in": taskGroupIDsToDelete},
			})
			if err != nil {
				log.Fatalf("Failed to delete task groups: %v", err)
			}
			log.Printf("Deleted %d task groups", deleteResult.DeletedCount)
		}
	}

	// Step 3: Get all tasks for the kept project
	var tasks []models.Task
	cursor, err = tasksCollection.Find(ctx, bson.M{"project_id": projectToKeep.ID})
	if err != nil {
		log.Fatalf("Failed to find tasks: %v", err)
	}
	if err := cursor.All(ctx, &tasks); err != nil {
		log.Fatalf("Failed to decode tasks: %v", err)
	}

	var taskToKeep *models.Task
	if len(tasks) == 0 {
		log.Println("No tasks found for the project.")
	} else {
		// Keep the first task, delete the rest
		taskToKeep = &tasks[0]
		log.Printf("Keeping task: %s (ID: %s)", taskToKeep.Name, taskToKeep.ID.Hex())

		if len(tasks) > 1 {
			taskIDsToDelete := make([]primitive.ObjectID, 0, len(tasks)-1)
			for i := 1; i < len(tasks); i++ {
				taskIDsToDelete = append(taskIDsToDelete, tasks[i].ID)
			}

			deleteResult, err := tasksCollection.DeleteMany(ctx, bson.M{
				"_id": bson.M{"$in": taskIDsToDelete},
			})
			if err != nil {
				log.Fatalf("Failed to delete tasks: %v", err)
			}
			log.Printf("Deleted %d tasks", deleteResult.DeletedCount)
		}

		// If we kept a task group, update the kept task to belong to it
		if taskGroupToKeep != nil && taskToKeep.TaskGroupID == nil {
			log.Printf("Updating task to belong to task group %s", taskGroupToKeep.Name)
			_, err := tasksCollection.UpdateOne(ctx,
				bson.M{"_id": taskToKeep.ID},
				bson.M{"$set": bson.M{"task_group_id": taskGroupToKeep.ID}},
			)
			if err != nil {
				log.Printf("Warning: Failed to update task's task_group_id: %v", err)
			}
		}
	}

	// Step 4: Delete all executions (clean slate)
	deleteResult, err := executionsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Failed to delete executions: %v", err)
	} else {
		log.Printf("Deleted %d executions", deleteResult.DeletedCount)
	}

	// Step 5: Delete task groups and tasks that don't belong to the kept project
	deleteResult, err = taskGroupsCollection.DeleteMany(ctx, bson.M{
		"project_id": bson.M{"$ne": projectToKeep.ID},
	})
	if err != nil {
		log.Printf("Warning: Failed to delete orphaned task groups: %v", err)
	} else {
		log.Printf("Deleted %d orphaned task groups", deleteResult.DeletedCount)
	}

	deleteResult, err = tasksCollection.DeleteMany(ctx, bson.M{
		"project_id": bson.M{"$ne": projectToKeep.ID},
	})
	if err != nil {
		log.Printf("Warning: Failed to delete orphaned tasks: %v", err)
	} else {
		log.Printf("Deleted %d orphaned tasks", deleteResult.DeletedCount)
	}

	// Summary
	log.Println("\n=== Cleanup Summary ===")
	log.Printf("Projects: 1 (kept: %s)", projectToKeep.Name)
	if taskGroupToKeep != nil {
		log.Printf("Task Groups: 1 (kept: %s)", taskGroupToKeep.Name)
	} else {
		log.Println("Task Groups: 0")
	}
	if taskToKeep != nil {
		log.Printf("Tasks: 1 (kept: %s)", taskToKeep.Name)
	} else {
		log.Println("Tasks: 0")
	}
	log.Println("Executions: 0 (all deleted)")
	log.Println("\nCleanup completed successfully!")
}

