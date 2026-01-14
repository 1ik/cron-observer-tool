import { SetMetadata } from '@nestjs/common';

export const CRON_TASK_METADATA_KEY = 'cron-task';

export interface CronTaskMetadata {
  taskName: string;
}

export const CronTask = (taskName: string) => SetMetadata(CRON_TASK_METADATA_KEY, { taskName });

