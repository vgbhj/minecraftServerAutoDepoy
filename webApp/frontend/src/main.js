import './assets/main.css'
import 'primeicons/primeicons.css';
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios';

axios.interceptors.request.use(config => {
  const adminPassword = localStorage.getItem('adminPassword');
  if (adminPassword) {
    config.headers['X-Admin-Password'] = adminPassword;
  }
  return config;
});

const app = createApp(App)

app.use(router)

app.mount('#app')
