import { Body, Controller, Logger, Post } from '@nestjs/common';

interface ExecutionPayload {
  task_name: string;
  execution_id: string;
}

@Controller('api/execute')
export class ExecutionController {
  private readonly logger = new Logger(ExecutionController.name);

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
    
    return {
      success: true,
      message: 'Execution received',
      received_at: new Date().toISOString(),
    };
  }
}

