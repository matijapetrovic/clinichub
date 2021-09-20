import utils from '@/utils';

export default {
  addMedicalDoctor(doctor) {
    return utils.clinicService.post('v1/doctors', doctor);
  },
  getAllDoctors(date) {
    return utils.apiClient.get(`api/medical-doctor/getAllOnDate/${date}`);
  },
  getAllDoctorsForClinic(clinicId) {
    return utils.clinicService.get(`v1/clinics/${clinicId}/doctors`);
  },
  getDoctorsForDateTime(payload) {
    return utils.apiClient.get(`api/medical-doctor/${payload.date}/${payload.time}`);
  },
  getWorkindCalendar() {
    return utils.apiClient.get('api/medical-doctor/schedule');
  },
  getWorkindCalendarByDoctorId(id) {
    return utils.apiClient.get(`api/medical-doctor/schedule/:${id}`);
  },
  addLeaveRequest(credentials) {
    return utils.apiClient.post('api/medical-doctor/addLeaveRequest', credentials);
  },
  finishAppointment(appointment) {
    return utils.apiClient.post('api/finished_appointment/add', appointment);
  },
  getAppointmentScheduleItem(appointmentId) {
    return utils.apiClient.get(`api/medical-doctor/getScheduleItem/${appointmentId}`);
  },
  getPreviousPatients() {
    return utils.apiClient.get('api/medical-doctor/previous-patients');
  },
  deleteDoctor(doctorId) {
    return utils.apiClient.post(`api/medical-doctor/delete/${doctorId}`);
  },
  updateDoctor(doctorId, doctorUpdate) {
    return utils.clinicService.put(`v1/doctors/${doctorId}`, doctorUpdate);
  },
};
