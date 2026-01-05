import { Module } from '@nestjs/common';
import { ExecutionController } from './execution/execution.controller';

@Module({
  imports: [],
  controllers: [ExecutionController],
  providers: [],
})
export class AppModule {}

