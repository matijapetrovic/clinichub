import utils from '@/utils';

export default {
  updateProfile(payload) {
    return utils.authService.post('profile', payload);
  },
  fetchProfile() {
    return utils.authService.get('profile');
  },
};
