// eslint-disable-next-line camelcase
import jwt_decode from 'jwt-decode';


import utils from '@/utils';
import api from './api';

export default {
  namespaced: true,
  state: {
    user: null,
    success: false,
    account: null,
  },
  mutations: {
    SET_USER_DATA(state, token) {
      const decoded = jwt_decode(token);
      console.log(decoded);
      state.user = decoded;
      localStorage.setItem('token', token);
      api.setToken(token);
    },
    CLEAR_USER_DATA() {
      localStorage.removeItem('token');
      window.location.reload();
    },
    SET_REGISTER_SUCCESS(state, success) {
      state.success = success;
    },
  },
  actions: {
    register({ commit, dispatch }, payload) {
      return api.register(payload)
        .then(() => {
          commit('SET_REGISTER_SUCCESS', true);
          const message = 'Registration request successful';
          dispatch('notifications/add', utils.successNotification(message), { root: true });
        })
        .catch((err) => {
          commit('SET_REGISTER_SUCCESS', false);
          dispatch('notifications/add', utils.errorNotification(err), { root: true });
        });
    },
    login({ state, commit, dispatch }, credentials) {
      return api.login(credentials)
        .then((response) => {
          commit('SET_USER_DATA', response.data.access_token);
          commit('SET_LOGGED_IN', state.user, { root: true });
          const message = 'Log in successful';
          dispatch('notifications/add', utils.successNotification(message), { root: true });
        })
        .catch((err) => {
          dispatch('notifications/add', utils.errorNotification(err), { root: true });
        });
    },
    logout({ commit, dispatch }) {
      commit('CLEAR_USER_DATA');
      const message = 'Log out successful';
      dispatch('notifications/add', utils.successNotification(message), { root: true });
    },
    changePassword({ dispatch }, credentials) {
      return api
        .changePassword(credentials)
        .then(() => {
          const message = 'Password change successful';
          dispatch('notifications/add', utils.successNotification(message), { root: true });
        })
        .catch((err) => {
          dispatch('notifications/add', utils.errorNotification(err), { root: true });
          throw err;
        });
    },
    activateAccount({ dispatch }, accountId) {
      return api.activateAccount(accountId)
        .then(() => {
          const message = 'Account activated successfully';
          dispatch('notifications/add', utils.successNotification(message), { root: true });
        })
        .catch((err) => {
          dispatch('notifications/add', utils.errorNotification(err), { root: true });
          throw err;
        });
    },
  },
  getters: {
    loggedIn(state) {
      return !!state.user;
    },
    passwordChanged(state) {
      return state.user && state.user.passwordChanged;
    },
    user(state) {
      return state.user;
    },
  },
};
