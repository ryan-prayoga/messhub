module.exports = {
  apps: [
    {
      name: 'messhub-frontend',
      cwd: __dirname,
      script: 'npm',
      args: 'run start',
      env: {
        NODE_ENV: 'production',
        HOST: '0.0.0.0',
        PORT: '4101',
        ORIGIN: 'http://127.0.0.1:4101',
        PUBLIC_APP_NAME: 'MessHub',
        PUBLIC_API_BASE_URL: '/api/v1'
      }
    }
  ]
};
