import utils from '@/utils';

export default {
  addClinicReview(clinicId, payload) {
    return utils.ratingService.post(`v1/clinics/${clinicId}/ratings`, payload);
  },
  addDoctorReview(doctorId, payload) {
    return utils.ratingService.post(`v1/doctors/${doctorId}/ratings`, payload);
  },
  fetchClinicsForReview() {
    return utils.ratingService.get('v1/clinics/to-rate');
  },
  fetchDoctorsForReview() {
    return utils.ratingService.get('v1/doctors/to-rate');
  },
};
