import utils from '@/utils';

export default {
  addClinic(clinic) {
    return utils.clinicService.post('v1/clinics', clinic);
  },
  getClinicNames() {
    return utils.apiClient.get('api/clinic/names');
  },
  fetchClinics(searchParams) {
    return utils.clinicService.get('v1/clinics', { params: searchParams });
  },
  getCurrentClinic() {
    return utils.apiClient.get('/api/clinic/getCurrent');
  },
  updateClinic(clinicId, clinic) {
    return utils.clinicService.put(`v1/clinics/${clinicId}`, clinic);
  },
  fetchPrices(clinicId) {
    return utils.clinicService.get(`v1/clinics/${clinicId}/prices`);
  },
  addPrice(clinicId, price) {
    return utils.clinicService.post(`v1/clinics/${clinicId}/prices`, price);
  },
  updatePrice(clinicId, price) {
    return utils.clinicService.put(`v1/clinics/${clinicId}/prices`, price);
  },
  fetchClinicProfile(clinicId) {
    return utils.clinicService.get(`v1/clinics/${clinicId}`);
  },
};
