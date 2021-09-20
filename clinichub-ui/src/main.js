import Vue from 'vue';
import App from './app/App.vue';
import router from './router';
import store from './store';
import vuetify from './plugins/vuetify';

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  vuetify,
  render: (h) => h(App),
  created() {
    const token = localStorage.getItem('token');
    if (token) {
      this.$store.commit('auth/SET_USER_DATA', token);
      this.$store.commit('SET_LOGGED_IN', this.$store.state.auth.user);
    }
  },
}).$mount('#app');
