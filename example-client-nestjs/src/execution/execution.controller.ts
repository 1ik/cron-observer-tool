import { Body, Controller, HttpException, HttpStatus, Logger, Post } from '@nestjs/common';
import { CronObserverService } from '../cron-observer/cron-observer.service';
import { CronTaskRouterService } from '../cron-observer/cron-task-router.service';
import { ExecutionContext } from '../cron-observer/interfaces';

interface ExecutionPayload {
  task_name: string;
  execution_id: string;
}

@Controller('api/execute')
export class ExecutionController {
  private readonly logger = new Logger(ExecutionController.name);

  constructor(
    private readonly routerService: CronTaskRouterService,
    private readonly cronObserver: CronObserverService,
  ) {}

  @Post()
  async handleExecution(@Body() body: ExecutionPayload) {
    this.logger.log('═══════════════════════════════════════════════════');
    this.logger.log('📨 Received execution request from cron-observer');
    this.logger.log('═══════════════════════════════════════════════════');
    this.logger.log(`Task Name: ${body.task_name}`);
    this.logger.log(`Execution ID: ${body.execution_id}`);
    this.logger.log(`Full Body: ${JSON.stringify(body, null, 2)}`);
    this.logger.log(`Timestamp: ${new Date().toISOString()}`);
    this.logger.log('═══════════════════════════════════════════════════');

    if (!body.task_name || !body.execution_id) {
      throw new HttpException(
        {
          success: false,
          error: 'Missing required fields: task_name and execution_id are required',
        },
        HttpStatus.BAD_REQUEST,
      );
    }

    const context: ExecutionContext = {
      executionId: body.execution_id,
      taskName: body.task_name,
    };

    try {
      const result = await this.routerService.executeTask(body.task_name, context);
      
      // Note: The handler should call cronObserver.success() or cronObserver.fail()
      // If it doesn't, the execution will remain PENDING
      // This is intentional - handlers control their own status updates
      
      return {
        success: true,
        message: 'Execution completed',
        result,
        received_at: new Date().toISOString(),
      };
    } catch (error) {
      this.logger.error(`Error executing task ${body.task_name}: ${error.message}`, error.stack);
      
      // Mark execution as failed if handler throws an error
      // (unless the handler already called fail() itself)
      try {
        await this.cronObserver.fail(
          body.execution_id,
          error instanceof Error ? error.message : 'Task execution failed',
        );
      } catch (failError) {
        this.logger.error(`Failed to mark execution as failed: ${failError.message}`);
      }
      
      throw new HttpException(
        {
          success: false,
          error: error.message || 'Task execution failed',
          task_name: body.task_name,
          execution_id: body.execution_id,
        },
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
}

