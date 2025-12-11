module.exports = {
  apps: [
    {
      name: 'janus-frontend',
      script: 'npm',
      args: 'run start',
      cwd: '/path/to/your/production/frontend', // 生产服务器上的项目路径
      instances: '1',
      exec_mode: 'cluster',
      env: {
        NODE_ENV: 'production',
      },
      env_production: {
        NODE_ENV: 'production',
        BACKEND_URL: 'http://localhost:10086',
      },
    },
  ],
};
