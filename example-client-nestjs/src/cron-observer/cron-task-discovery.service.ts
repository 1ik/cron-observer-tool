import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { ModuleRef } from '@nestjs/core';
import { CronTaskRouterService } from './cron-task-router.service';
import { CRON_TASK_METADATA_KEY } from './cron-task.decorator';

@Injectable()
export class CronTaskDiscoveryService implements OnModuleInit {
  private readonly logger = new Logger(CronTaskDiscoveryService.name);

  constructor(
    private moduleRef: ModuleRef,
    private routerService: CronTaskRouterService,
  ) {}

  onModuleInit() {
    this.discoverHandlers();
  }

  private discoverHandlers(): void {
    let handlerCount = 0;

    // Try to discover from all registered providers
    try {
      const container = (this.moduleRef as any).container;
      if (container && container.modules) {
        for (const module of container.modules.values()) {
          const providers = module.providers || new Map();
          
          for (const [token] of providers) {
            if (typeof token === 'function') {
              try {
                const instance = this.moduleRef.get(token, { strict: false });
                if (instance) {
                  const discovered = this.discoverHandlersInInstance(instance);
                  handlerCount += discovered;
                }
              } catch (error) {
                // Skip if provider is not available
              }
            }
          }
        }
      }
    } catch (error) {
      this.logger.debug(`Container-based discovery failed: ${error.message}`);
    }

    this.logger.log(`Discovered ${handlerCount} cron task handler(s)`);
  }

  private discoverHandlersInInstance(instance: any): number {
    let count = 0;
    const prototype = Object.getPrototypeOf(instance);
    const methodNames = Object.getOwnPropertyNames(prototype).filter(
      (name) => name !== 'constructor' && typeof prototype[name] === 'function',
    );

    for (const methodName of methodNames) {
      const metadata: { taskName: string } | undefined = Reflect.getMetadata(
        CRON_TASK_METADATA_KEY,
        prototype[methodName],
      );

      if (metadata) {
        this.routerService.registerHandler(metadata.taskName, instance, methodName);
        count++;
      }
    }

    return count;
  }
}

