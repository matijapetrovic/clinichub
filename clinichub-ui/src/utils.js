import axios from 'axios';

export default {
  apiClient: axios.create({
    baseURL: process.env.VUE_APP_API_URL,
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    timeout: 10000,
  }),
  clinicService: axios.create({
    baseURL: 'http://localhost:8081',
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    timeout: 10000,
  }),
  ratingService: axios.create({
    baseURL: 'http://localhost:8082',
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    timeout: 10000,
  }),
  schedulingService: axios.create({
    baseURL: 'http://localhost:8083',
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    timeout: 10000,
  }),
  authService: axios.create({
    baseURL: 'http://127.0.0.1:5000',
    withCredentials: false,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    timeout: 10000,
  }),
  successNotification(msg) {
    return {
      text: msg,
      color: 'success',
    };
  },
  errorNotification(err) {
    console.log(err.response);
    let text = err.response ? err.response.data.message : err;
    if (!text) {
      text = err.response.status;
      if (!text || text === 401) {
        text = 'Invalid username/password';
      }
    }
    return {
      text,
      color: 'error',
    };
  },
  roleRootPath(role) {
    switch (role) {
      case 'patient':
        return '/patient';
      case 'admin':
        return '/clinic-center-admin';
      default:
        return '';
    }
  },
};
