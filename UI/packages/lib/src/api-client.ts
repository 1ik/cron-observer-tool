import { makeApi, Zodios, type ZodiosOptions } from "@zodios/core";
import { z } from "zod";

const models_Project = z
  .object({
    api_key: z.string(),
    created_at: z.string(),
    description: z.string(),
    execution_endpoint: z.string(),
    id: z.string(),
    name: z.string(),
    updated_at: z.string(),
    uuid: z.string(),
  })
  .partial()
  .passthrough();
const models_ErrorResponse = z
  .object({ details: z.array(z.string()), error: z.string() })
  .partial()
  .passthrough();
const models_CreateProjectRequest = z
  .object({
    description: z.string().max(1000).optional(),
    execution_endpoint: z.string().optional(),
    name: z.string().min(1).max(255),
  })
  .passthrough();
const models_TaskGroupStatus = z.enum(["ACTIVE", "PAUSED", "DISABLED"]);
const models_TaskGroup = z
  .object({
    created_at: z.string(),
    description: z.string(),
    end_time: z.string(),
    id: z.string(),
    name: z.string(),
    project_id: z.string(),
    start_time: z.string(),
    status: models_TaskGroupStatus,
    timezone: z.string(),
    updated_at: z.string(),
    uuid: z.string(),
  })
  .partial()
  .passthrough();
const models_CreateTaskGroupRequest = z
  .object({
    description: z.string().max(1000).optional(),
    end_time: z.string().optional(),
    name: z.string().min(1).max(255),
    project_id: z.string(),
    start_time: z.string().optional(),
    status: models_TaskGroupStatus.optional(),
    timezone: z.string().optional(),
  })
  .passthrough();
const models_UpdateTaskGroupRequest = z
  .object({
    description: z.string().max(1000).optional(),
    end_time: z.string().optional(),
    name: z.string().min(1).max(255),
    start_time: z.string().optional(),
    status: models_TaskGroupStatus.optional(),
    timezone: z.string().optional(),
  })
  .passthrough();
const models_FrequencyUnit = z.enum(["s", "m", "h"]);
const models_Frequency = z
  .object({ unit: models_FrequencyUnit, value: z.number().int().gte(1) })
  .passthrough();
const models_TimeRange = z
  .object({ end: z.string(), frequency: models_Frequency, start: z.string() })
  .passthrough();
const models_ScheduleConfig = z
  .object({
    cron_expression: z.string().optional(),
    days_of_week: z.array(z.number().int()).optional(),
    exclusions: z.array(z.number().int()).optional(),
    time_range: models_TimeRange.optional(),
    timezone: z.string(),
  })
  .passthrough();
const models_ScheduleType = z.enum(["RECURRING", "ONEOFF"]);
const models_TaskStatus = z.enum(["ACTIVE", "PAUSED", "DISABLED"]);
const models_HTTPTriggerConfig = z
  .object({
    body: z.unknown().optional(),
    headers: z.record(z.string()).optional(),
    method: z.string(),
    timeout: z.number().int().gte(1).lte(300).optional(),
    url: z.string(),
  })
  .passthrough();
const models_TriggerType = z.literal("HTTP");
const models_TriggerConfig = z
  .object({ http: models_HTTPTriggerConfig, type: models_TriggerType })
  .partial()
  .passthrough();
const models_Task = z
  .object({
    created_at: z.string(),
    description: z.string(),
    id: z.string(),
    metadata: z.object({}).partial().passthrough(),
    name: z.string(),
    project_id: z.string(),
    schedule_config: models_ScheduleConfig,
    schedule_type: models_ScheduleType,
    status: models_TaskStatus,
    task_group_id: z.string(),
    trigger_config: models_TriggerConfig,
    updated_at: z.string(),
    uuid: z.string(),
  })
  .partial()
  .passthrough();
const models_CreateTaskRequest = z
  .object({
    description: z.string().max(1000).optional(),
    metadata: z.object({}).partial().passthrough().optional(),
    name: z.string().min(1).max(255),
    project_id: z.string(),
    schedule_config: models_ScheduleConfig,
    schedule_type: models_ScheduleType,
    status: models_TaskStatus.optional(),
    task_group_id: z.string().optional(),
  })
  .passthrough();
const models_UpdateTaskRequest = z
  .object({
    description: z.string().max(1000).optional(),
    metadata: z.object({}).partial().passthrough().optional(),
    name: z.string().min(1).max(255),
    schedule_config: models_ScheduleConfig,
    schedule_type: models_ScheduleType,
    status: models_TaskStatus.optional(),
    task_group_id: z.string().optional(),
  })
  .passthrough();

export const schemas = {
  models_Project,
  models_ErrorResponse,
  models_CreateProjectRequest,
  models_TaskGroupStatus,
  models_TaskGroup,
  models_CreateTaskGroupRequest,
  models_UpdateTaskGroupRequest,
  models_FrequencyUnit,
  models_Frequency,
  models_TimeRange,
  models_ScheduleConfig,
  models_ScheduleType,
  models_TaskStatus,
  models_HTTPTriggerConfig,
  models_TriggerType,
  models_TriggerConfig,
  models_Task,
  models_CreateTaskRequest,
  models_UpdateTaskRequest,
};

const endpoints = makeApi([
  {
    method: "get",
    path: "/projects",
    alias: "getProjects",
    description: `Retrieve a list of all projects`,
    requestFormat: "json",
    response: z.array(models_Project),
    errors: [
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/projects",
    alias: "postProjects",
    description: `Create a new project with auto-generated UUID and API key`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Project creation request`,
        type: "Body",
        schema: models_CreateProjectRequest,
      },
    ],
    response: models_Project,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/projects/:project_id/task-groups",
    alias: "getProjectsProject_idtaskGroups",
    description: `Retrieve all task groups belonging to a project`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.array(models_TaskGroup),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/projects/:project_id/task-groups",
    alias: "postProjectsProject_idtaskGroups",
    description: `Create a new task group in a project`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Task group creation request`,
        type: "Body",
        schema: models_CreateTaskGroupRequest,
      },
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: models_TaskGroup,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/projects/:project_id/task-groups/:group_uuid",
    alias: "getProjectsProject_idtaskGroupsGroup_uuid",
    description: `Retrieve a task group by its UUID`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: models_TaskGroup,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 404,
        description: `Not Found`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/projects/:project_id/task-groups/:group_uuid",
    alias: "putProjectsProject_idtaskGroupsGroup_uuid",
    description: `Update an existing task group`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Task group update request`,
        type: "Body",
        schema: models_UpdateTaskGroupRequest,
      },
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: models_TaskGroup,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 404,
        description: `Not Found`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/projects/:project_id/task-groups/:group_uuid",
    alias: "deleteProjectsProject_idtaskGroupsGroup_uuid",
    description: `Delete an existing task group`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/projects/:project_id/task-groups/:group_uuid/start",
    alias: "postProjectsProject_idtaskGroupsGroup_uuidstart",
    description: `Manually start all tasks in a task group`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.record(z.string()),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/projects/:project_id/task-groups/:group_uuid/stop",
    alias: "postProjectsProject_idtaskGroupsGroup_uuidstop",
    description: `Manually stop all tasks in a task group`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.record(z.string()),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/projects/:project_id/task-groups/:group_uuid/tasks",
    alias: "getProjectsProject_idtaskGroupsGroup_uuidtasks",
    description: `Retrieve all tasks belonging to a task group`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "group_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.array(models_Task),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 404,
        description: `Not Found`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/projects/:project_id/tasks",
    alias: "getProjectsProject_idtasks",
    description: `Retrieve all tasks belonging to a project`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.array(models_Task),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/projects/:project_id/tasks",
    alias: "postProjectsProject_idtasks",
    description: `Create a new scheduled task in a project`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Task creation request`,
        type: "Body",
        schema: models_CreateTaskRequest,
      },
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: models_Task,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/projects/:project_id/tasks/:task_uuid",
    alias: "putProjectsProject_idtasksTask_uuid",
    description: `Update an existing scheduled task`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Task update request`,
        type: "Body",
        schema: models_UpdateTaskRequest,
      },
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "task_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: models_Task,
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 404,
        description: `Not Found`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/projects/:project_id/tasks/:task_uuid",
    alias: "deleteProjectsProject_idtasksTask_uuid",
    description: `Delete an existing scheduled task`,
    requestFormat: "json",
    parameters: [
      {
        name: "project_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "task_uuid",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad Request`,
        schema: models_ErrorResponse,
      },
      {
        status: 500,
        description: `Internal Server Error`,
        schema: models_ErrorResponse,
      },
    ],
  },
]);

export const api = new Zodios(endpoints);

export function createApiClient(baseUrl: string, options?: ZodiosOptions) {
  return new Zodios(baseUrl, endpoints, options);
}
