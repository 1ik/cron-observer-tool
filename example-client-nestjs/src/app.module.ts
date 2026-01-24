import { Module } from '@nestjs/common';
import { CronObserverModule } from './cron-observer/cron-observer.module';
import { ExecutionController } from './execution/execution.controller';
import { TaskHandlersService } from './execution/task-handlers.service';

@Module({
  imports: [
    CronObserverModule.forRoot({
      backendUrl: process.env.CRON_OBSERVER_URL || 'http://localhost:8080',
      apiKey: process.env.CRON_OBSERVER_API_KEY || '',
    }),
  ],
  controllers: [ExecutionController],
  providers: [TaskHandlersService],
})
export class AppModule {}

