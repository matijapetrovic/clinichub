import ClinicCenterAdmin from '@/views/clinic-center-admin/ClinicCenterAdmin.vue';
import ClinicCenterAdminHome from '@/views/clinic-center-admin/ClinicCenterAdminHome.vue';
import AddClinic from '@/views/clinic-center-admin/AddClinic.vue';
import AddAppointmentTypeForm from '@/app/appointment_type/components/AddAppointmentTypeForm.vue';
import AppointmentTypeTable from '@/app/appointment_type/components/AppointmentTypeTable.vue';
import AddMedicalDoctorForm from '@/app/medical_doctor/components/AddMedicalDoctorForm.vue';
import AllDoctorsView from '@/views/medical_doctor/AllDoctorsView.vue';
import AllClinics from '@/views/clinic-center-admin/AllClinics.vue';
import PriceListTable from '@/views/clinic/PriceListTable.vue';

export default {
  path: '/clinic-center-admin',
  name: 'clinic-center-admin',
  component: ClinicCenterAdmin,
  meta: {
    requiresAuth: true,
    role: 'admin',
  },
  children: [
    {
      path: '',
      component: ClinicCenterAdminHome,
    },
    {
      path: '/add-clinic',
      component: AddClinic,
    },
    {
      path: '/clinics',
      component: AllClinics,
    },
    {
      path: '/add-appointment-type',
      component: AddAppointmentTypeForm,
    },
    {
      path: '/appointment-types',
      component: AppointmentTypeTable,
    },
    {
      path: '/appointment-prices',
      component: PriceListTable,
    },
    {
      path: '/add-doctor',
      component: AddMedicalDoctorForm,
    },
    {
      path: '/doctors',
      component: AllDoctorsView,
    },
  ],
};
