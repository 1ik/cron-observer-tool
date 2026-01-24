import { DynamicModule, Module, Provider } from '@nestjs/common';
import { CronObserverService } from './cron-observer.service';
import { CronTaskDiscoveryService } from './cron-task-discovery.service';
import { CronTaskRouterService } from './cron-task-router.service';
import { CronObserverConfig } from './interfaces';

@Module({})
export class CronObserverModule {
  static forRoot(config: CronObserverConfig): DynamicModule {
    const cronObserverServiceProvider: Provider = {
      provide: CronObserverService,
      useValue: new CronObserverService(config.backendUrl, config.apiKey),
    };

    return {
      module: CronObserverModule,
      providers: [
        cronObserverServiceProvider,
        CronTaskRouterService,
        CronTaskDiscoveryService,
      ],
      exports: [CronObserverService, CronTaskRouterService],
    };
  }
}

