import utils from '@/utils';

export default {
  addMedicalDoctor(credentials) {
    return utils.apiClient.post('api/medical-doctor/add', credentials);
  },
  fetchDoctors(clinicId, searchParams) {
    return utils.clinicService.get(`v1/clinics/${clinicId}/doctors`, { params: searchParams });
  },
  fetchAvailableTimes(doctorId, date) {
    return utils.apiClient.get(`api/medical-doctor/${doctorId}/available_times/${date}`);
  },
  fetchAvailableTimesWithoutState(doctorId, date) {
    return utils.apiClient.get(`api/medical-doctor/${doctorId}/available_times/${date}`);
  },
};
