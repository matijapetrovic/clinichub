import utils from '@/utils';

export default {
  setToken(token) {
    utils.authService.defaults.headers.common.Authorization = `Bearer ${token}`;
    utils.clinicService.defaults.headers.common.Authorization = `Bearer ${token}`;
    utils.schedulingService.defaults.headers.common.Authorization = `Bearer ${token}`;
    utils.ratingService.defaults.headers.common.Authorization = `Bearer ${token}`;
  },
  register(payload) {
    return utils.authService.post('register', payload);
  },
  login(credentials) {
    const form = new FormData();
    form.append('username', credentials.username);
    form.append('password', credentials.password);
    form.append('grant_type', 'password');
    return utils.authService.post('oauth/token', form, {
      auth: {
        username: 'o8JNZJ2F8wDm92ca2MNU5S4l',
        password: 'Cp4LeK2RX9hzKebRI2fI6t9dKG5WiUSug0xeX7ZMwMZ8gKaT',
      },
    });
  },
  changePassword(credentials) {
    return utils.authService.post('change-password', credentials);
  },
  activateAccount(accountId) {
    return utils.apiClient.post(`api/auth/activate/${accountId}`);
  },
};
