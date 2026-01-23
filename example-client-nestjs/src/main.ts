import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Enable CORS
  app.enableCors();
  
  const port = process.env.PORT || 5202;
  await app.listen(port);
  
  console.log(`🚀 NestJS Example Client running on http://localhost:${port}`);
  console.log(`📡 Execution endpoint: http://localhost:${port}/api/execute`);
}

bootstrap();

