<template>
  <BaseLayoutSidebar>
    <template v-slot:sidebar>
      <SideNav
        :items="navItems"
      />
    </template>
    <v-content>
      <router-view></router-view>
    </v-content>
  </BaseLayoutSidebar>
</template>

<script>
import BaseLayoutSidebar from '@/app/shared/_components/_layouts/BaseLayoutSidebar.vue';
import SideNav from '@/app/shared/_components/_navigation/SideNav.vue';

export default {
  name: 'ClinicAdmin',
  components: {
    SideNav,
    BaseLayoutSidebar,
  },
  data: () => ({
    navItems: [
      {
        text: 'Clinics',
        icon: 'mdi-hospital-box-outline',
        children: [
          { text: 'Add', path: '/add-clinic', icon: 'add' },
          { text: 'All clinics', path: '/clinics', icon: 'mdi-hospital-box' },
        ],
      },
      {
        text: 'Appointment Types',
        icon: 'mdi-clipboard-pulse-outline',
        children: [
          { text: 'Add', path: '/add-appointment-type', icon: 'add' },
          { text: 'Search', path: '/appointment-types', icon: 'search' },
        ],
      },
      {
        text: 'Doctors',
        path: '/doctors',
        icon: 'mdi-doctor',
        children: [
          { text: 'Add', path: '/add-doctor', icon: 'add' },
          { text: 'All doctors', path: '/doctors', icon: 'mdi-account-multiple-outline' },
        ],
      },
      {
        text: 'Appointment Prices',
        icon: 'mdi-book-open',
        path: '/appointment-prices',
      },
    ],
  }),
};
</script>
