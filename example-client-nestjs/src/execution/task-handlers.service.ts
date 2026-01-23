import { Injectable } from '@nestjs/common';
import { CronObserverService } from '../cron-observer/cron-observer.service';
import { CronTask } from '../cron-observer/cron-task.decorator';
import { ExecutionContext } from '../cron-observer/interfaces';

@Injectable()
export class TaskHandlersService {
  constructor(private readonly cronObserver: CronObserverService) {}

  @CronTask('daily-report')
  async handleDailyReport(ctx: ExecutionContext) {
    await this.cronObserver.log(ctx.executionId, 'Starting daily report...', 'info');
    
    try {
      // Simulate some work
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      await this.cronObserver.log(ctx.executionId, 'Processing data...', 'info');
      await new Promise(resolve => setTimeout(resolve, 500));
      
      await this.cronObserver.log(ctx.executionId, 'Report complete', 'info');
      
      // Acknowledge successful completion
      await this.cronObserver.success(ctx.executionId);
      
      return {
        success: true,
        message: 'Daily report processed',
        executionId: ctx.executionId,
      };
    } catch (error) {
      await this.cronObserver.fail(ctx.executionId, error instanceof Error ? error.message : 'Unknown error');
      throw error;
    }
  }

  @CronTask('data-sync')
  async handleDataSync(ctx: ExecutionContext) {
    await this.cronObserver.log(ctx.executionId, 'Starting data synchronization...', 'info');
    
    try {
      await this.cronObserver.log(ctx.executionId, 'Connecting to database...', 'info');
      await new Promise(resolve => setTimeout(resolve, 500));
      
      await this.cronObserver.log(ctx.executionId, 'Syncing records...', 'info');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      await this.cronObserver.log(ctx.executionId, 'Data sync completed successfully', 'info');
      
      // Acknowledge successful completion
      await this.cronObserver.success(ctx.executionId);
      
      return {
        success: true,
        message: 'Data synchronized',
        executionId: ctx.executionId,
      };
    } catch (error) {
      await this.cronObserver.log(ctx.executionId, `Error: ${error.message}`, 'error');
      await this.cronObserver.fail(ctx.executionId, error instanceof Error ? error.message : 'Unknown error');
      throw error;
    }
  }

  @CronTask('Valid Task')
  async handleValidTask(ctx: ExecutionContext) {
    await this.cronObserver.log(ctx.executionId, 'Starting Valid Task execution...', 'info');
    
    try {
      // Simulate 2 seconds of execution
      await this.cronObserver.log(ctx.executionId, 'Processing Valid Task...', 'info');
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      await this.cronObserver.log(ctx.executionId, 'Valid Task completed successfully', 'info');
      
      // Acknowledge successful completion
      await this.cronObserver.success(ctx.executionId);
      
      return {
        success: true,
        message: 'Valid Task executed successfully',
        executionId: ctx.executionId,
      };
    } catch (error) {
      await this.cronObserver.log(ctx.executionId, `Error: ${error.message}`, 'error');
      await this.cronObserver.fail(ctx.executionId, error instanceof Error ? error.message : 'Unknown error');
      throw error;
    }
  }
}

