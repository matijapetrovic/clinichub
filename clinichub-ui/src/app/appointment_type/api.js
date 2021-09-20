import utils from '@/utils';

export default {
  addAppointmentType(appointmentType) {
    return utils.clinicService.post('v1/appointment-types', appointmentType);
  },
  fetchAppointmentTypes() {
    return utils.clinicService.get('v1/appointment-types');
  },
  removeAppointmentType(id) {
    return utils.apiClient.post('api/appointment-type/delete', id);
  },
  changeAppointmentType(appointmentTypeId, appointmentType) {
    return utils.clinicService.put(`v1/appointment-types/${appointmentTypeId}`, appointmentType);
  },
};
