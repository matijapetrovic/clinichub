import utils from '@/utils';

export default {
  scheduleAppointment(appointmentRequest) {
    return utils.schedulingService.post('v1/appointments', appointmentRequest);
  },
  scheduleDoctorsAppointment(appointmentRequest) {
    return utils.apiClient.post('api/appointment-request/addForDoctor', appointmentRequest);
  },
  addAppointment(appointment) {
    return utils.apiClient.post('api/appointment/add', appointment);
  },
  fetchScheduledAppointments() {
    return utils.schedulingService.get('v1/appointments');
  },
  getCurrentAppointment(patientId) {
    return utils.apiClient.get(`api/appointment/getCurrent/${patientId}`);
  },
  cancelScheduledAppointment(appointmentId) {
    return utils.apiClient.post(`/api/appointment/${appointmentId}/cancel`);
  },
};
