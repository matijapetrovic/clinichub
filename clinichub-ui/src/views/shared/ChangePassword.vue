<template>
  <div class="change-password">
    <v-row
      align="center"
      justify="center"
    >
      <v-col
        cols="12"
        sm="8"
        md="4"
      >
        <ChangePasswordForm/>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import ChangePasswordForm from '@/app/authentication/_components/ChangePasswordForm.vue';
import { mapGetters } from 'vuex';

export default {
  name: 'ChangePassword',
  components: {
    ChangePasswordForm,
  },
  computed: {
    ...mapGetters('auth', ['passwordChanged']),
  },
};
</script>
