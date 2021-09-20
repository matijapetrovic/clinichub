<template>
  <div>
    <v-data-table
      :headers="headers"
      :items="items"
      :items-per-page="5"
      class="elevation-1"
    >
      <template v-slot:body="{ items }">
          <tbody>
            <tr v-for="item in items" :key="item.id">
              <td>{{ item.doctorFullName }}</td>
              <td>{{ item.time.split("T")[0] }}</td>
              <td>{{ item.time.split("T")[1].substring(0, 5) }}</td>
              <!-- <td>
                <v-btn :disabled="!canCancel(item)" @click="openCancelDialog(item)">
                Cancel</v-btn>
              </td> -->
            </tr>
          </tbody>
        </template>
    </v-data-table>
    <v-dialog
      v-model="dialog"
      max-width="300"
    >
      <CancelAppointmentDialog
        @ok="cancel"
        @cancel="closeDialog"
      />
    </v-dialog>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import CancelAppointmentDialog from './CancelAppointmentDialog.vue';

export default {
  name: 'ScheduledAppointments',
  components: {
    CancelAppointmentDialog,
  },
  data: () => ({
    headers: [
      { text: 'Doctor', value: 'doctorFullName', align: 'start' },
      { text: 'Date', value: 'time' },
      { text: 'Time', value: 'time' },
    ],
    dialog: false,
  }),
  props: {
    items: {
      type: Array,
      required: true,
    },
  },
  computed: {
  },
  methods: {
    ...mapActions('appointment', ['cancelScheduledAppointment']),
    cancel() {
      this.cancelScheduledAppointment(this.appointment.id)
        .then(() => {
          this.closeDialog();
        });
    },
    openCancelDialog(appointment) {
      this.appointment = appointment;
      this.dialog = true;
    },
    closeDialog() {
      this.dialog = false;
    },
    canCancel(appointment) {
      function addDays(date, days) {
        const res = new Date(date);
        res.setDate(res.getDate() + days);
        return res;
      }
      const dateStr = `${appointment.date}T${appointment.time}`;
      const appDate = Date.parse(dateStr);
      return addDays(Date.now(), 1).getTime() < new Date(appDate).getTime();
    },
  },
};
</script>
