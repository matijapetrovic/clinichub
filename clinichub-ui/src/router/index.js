import Vue from 'vue';
import VueRouter from 'vue-router';
// eslint-disable-next-line camelcase
import jwt_decode from 'jwt-decode';
import utils from '@/utils';
import Home from '@/views/Home.vue';

import clinicCenterAdmin from './routes/clinic-center-admin';
import unregisteredUser from './routes/unregistered';
import registeredUser from './routes/user';
import patient from './routes/patient';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: {
      requiresAuth: true,
    },
  },
  clinicCenterAdmin,
  unregisteredUser,
  registeredUser,
  patient,
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!token) {
      next('/login');
    } else {
      const user = jwt_decode(token);
      console.log(user);
      if (to.matched.some((record) => record.path === '')) {
        next(utils.roleRootPath(user.role));
      } else if (to.matched.some((record) => record.meta.role
                                              && user.role !== record.meta.role)) {
        next(utils.roleRootPath(user.role));
      } else {
        next();
      }
    }
  } else if (token) {
    next('/');
  } else {
    next();
  }
});

export default router;
