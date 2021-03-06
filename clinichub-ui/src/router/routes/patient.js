import Patient from '@/views/patient/Patient.vue';
import PatientReviews from '@/views/patient/PatientReviews.vue';
import PatientSearchClinics from '@/views/patient/PatientSearchClinics.vue';
import PatientSearchDoctors from '@/views/patient/PatientSearchDoctors.vue';
import PatientScheduledAppointments from '@/views/patient/PatientScheduledAppointments.vue';
import PatientClinicProfile from '@/views/patient/PatientClinicProfile.vue';
import PatientHome from '@/views/patient/PatientHome.vue';

export default {
  path: '/patient',
  name: 'patient',
  component: Patient,
  meta: {
    requiresAuth: true,
    role: 'patient',
  },
  children: [
    {
      path: '/search-clinics',
      component: PatientSearchClinics,
    },
    {
      path: '/search-doctors/:clinic_id',
      component: PatientSearchDoctors,
    },
    {
      path: '/scheduled-appointments',
      component: PatientScheduledAppointments,
    },
    {
      path: '',
      component: PatientHome,
    },
    {
      path: '/clinic/:clinic_id',
      component: PatientClinicProfile,
    },
    {
      path: '/reviews',
      component: PatientReviews,
    },
  ],
};
