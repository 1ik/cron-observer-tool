import { Injectable, Logger } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import { ExecutionContext } from './interfaces';

interface TaskHandler {
  instance: any;
  methodName: string;
  taskName: string;
}

@Injectable()
export class CronTaskRouterService {
  private readonly logger = new Logger(CronTaskRouterService.name);
  private handlers: Map<string, TaskHandler> = new Map();

  constructor(private moduleRef: ModuleRef) {}

  registerHandler(taskName: string, instance: any, methodName: string): void {
    this.handlers.set(taskName, {
      instance,
      methodName,
      taskName,
    });
    this.logger.log(`Registered handler for task: ${taskName} -> ${instance.constructor.name}.${methodName}`);
  }

  async executeTask(taskName: string, context: ExecutionContext): Promise<any> {
    const handler = this.handlers.get(taskName);
    if (!handler) {
      throw new Error(`No handler found for task: ${taskName}`);
    }

    this.logger.log(`Executing task: ${taskName} with execution ID: ${context.executionId}`);
    
    try {
      const method = handler.instance[handler.methodName];
      if (typeof method !== 'function') {
        throw new Error(`Method ${handler.methodName} is not a function`);
      }
      
      return await method.call(handler.instance, context);
    } catch (error) {
      this.logger.error(`Error executing task ${taskName}: ${error.message}`, error.stack);
      throw error;
    }
  }

  getRegisteredTasks(): string[] {
    return Array.from(this.handlers.keys());
  }
}

