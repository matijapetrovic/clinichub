<template>
  <div>
    <v-row>
      <v-col>
        <ClinicMap
          :coords="coords"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col md="4">
        <v-rating
          :value="clinic.rating.rating"
          color="amber"
          dense
          half-increments
          readonly
          class="float-left"
        >
        </v-rating>
        <span class="grey--text ml-4">
          {{ parseFloat(clinic.rating.rating).toFixed(2) }} ({{ clinic.rating.count }})
        </span>
      </v-col>
      <v-col md="8">
        <div>{{ clinic.description }}</div>
      </v-col>
    </v-row>
    <v-dialog
      v-model="dialog"
      max-width="500"
    >
      <ScheduleAppointmentDialog
        :appointment="appointment"
        @cancelled="closeConfirmScheduleDialog"
        @scheduled="schedule"
      />
    </v-dialog>
  </div>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import ScheduleAppointmentDialog from '@/app/appointment/_components/ScheduleAppointmentDialog.vue';
import ClinicMap from './ClinicMap.vue';

export default {
  name: 'ClinicProfile',
  components: {
    ClinicMap,
    ScheduleAppointmentDialog,
  },
  data: () => ({
    appointment: {},
    dialog: false,
  }),
  props: {
    clinic: {
      required: true,
      type: Object,
    },
    coords: {
      required: true,
      type: Array,
    },
  },
  methods: {
    ...mapActions('predefinedAppointment', ['schedulePredefinedAppointment']),
    openConfirmScheduleDialog(appointment) {
      this.appointment = appointment;
      this.dialog = true;
    },
    closeConfirmScheduleDialog() {
      this.dialog = false;
    },
    seeDoctors() {
      this.$router.push(`/search-doctors/${this.clinic.id}`);
    },
    schedule() {
      this.schedulePredefinedAppointment(this.appointment.id)
        .then(() => {
          this.closeConfirmScheduleDialog();
        })
        .catch((err) => {
          this.closeConfirmScheduleDialog();
          if (err.response.status === 404) {
            const location = this.$route.fullPath;
            this.$router.replace('/');
            this.$nextTick(() => this.$router.replace(location));
          }
        });
    },
  },
  computed: {
    ...mapState('predefinedAppointment', ['predefinedAppointments']),
  },
};
</script>

<style scoped>
.appointment-list {
  max-height: 300px;
  overflow-y: auto;
}
.appointment-list::-webkit-scrollbar
{
    width: 6px;
    background-color: #F5F5F5;
}

.appointment-list::-webkit-scrollbar-thumb
{
    background-color: #000000;
}
</style>
