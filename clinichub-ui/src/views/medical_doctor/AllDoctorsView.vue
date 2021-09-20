<template>
  <v-container>
    <v-row
      class="mb-6"
      no-gutters
    >

    <v-col cols="4">
      <v-select
        :items="clinics"
        v-model="clinic"
        label="Clinic"
        item-text="name"
        return-object
        prepend-icon="mdi-hospital-box-outline"
      ></v-select>
    </v-col>
    <v-col cols="3">
      <v-btn
       :disabled="!clinic"
        color="primary" class="mt-5 ml-5" @click="findDoctors">Find doctors</v-btn>
    </v-col>
    <v-col
      cols="12"
      lg="12"
    >
    <v-data-table
      :headers="headers"
      :items="doctors"
      :search="search"
      class="elevation-1"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>Medical Doctors</v-toolbar-title>
          <v-spacer auto></v-spacer>
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search"
            single-line
            hide-details
          ></v-text-field>
          <v-divider
            class="mx-4"
            inset
            vertical
          ></v-divider>
          <v-spacer></v-spacer>
          <v-dialog v-model="dialog" max-width="500px">
            <v-card>
              <v-card-title>
                <span class="headline">{{ formTitle }}</span>
              </v-card-title>

              <v-card-text>
                <v-form ref="form">
                  <v-container>
                    <v-row>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field v-model="editedItem.firstName" label="First name">
                        </v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          v-model="editedItem.lastName"
                          label="Last name">
                        </v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-select
                          :items="dayHours"
                          v-model="editedItem.workStart"
                          label="Work start"
                          :rules="[requiredRule]"
                        ></v-select>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-select
                          :items="dayHours"
                          v-model="editedItem.workEnd"
                          label="Work end"
                          :rules="[requiredRule]"
                        ></v-select>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-form>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
                <v-btn
                :disabled="saveDisabled"
                 color="blue darken-1"
                  text @click="save">Save</v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-toolbar>
      </template>
      <template v-slot:item.actions="{ item }">
        <v-icon
          small
          @click="editItem(item)"
        >
          mdi-pencil
        </v-icon>
      </template>
      <template v-slot:no-data>
        No doctors available
      </template>
    </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions, mapState } from 'vuex';

export default {
  data: () => ({
    dayHours: [
      '00:00', '01:00', '02:00', '03:00', '04:00', '05:00',
      '06:00', '07:00', '08:00', '09:00', '10:00', '11:00',
      '12:00', '13:00', '14:00', '15:00', '16:00', '17:00',
      '18:00', '19:00', '20:00', '21:00', '22:00', '23:00',
    ],
    dialog: false,
    search: '',
    headers: [
      {
        text: 'First name',
        align: 'start',
        value: 'firstName',
      },
      { text: 'Last name', value: 'lastName' },
      { text: 'Specialization', value: 'specialization.name' },
      { text: 'Start of work', value: 'workStart' },
      { text: 'End of work', value: 'workEnd' },
      { text: 'Actions', value: 'actions', sortable: false },
    ],
    editedIndex: -1,
    editedItem: {
      firstName: '',
      lastName: '',
      workStart: '',
      workEnd: '',
    },
    defaultItem: {
      firstName: '',
      lastName: '',
      workStart: '',
      workEnd: '',
    },
    clinic: null,
  }),

  computed: {
    ...mapState('medicalDoctor', ['doctors']),
    ...mapState('clinic', ['clinics']),
    formTitle() {
      return this.editedIndex === -1 ? 'New Item' : 'Edit Item';
    },
    saveDisabled() {
      return this.$refs.form && this.$refs.form.valid;
    },
    requiredRule() {
      return (value) => !!value || 'Required';
    },
  },

  watch: {
    dialog(val) {
      return () => val || this.close();
    },
    search(val) {
      return () => val && val !== this.select && this.querySelections(val);
    },
  },

  created() {
    this.initialize();
  },

  methods: {
    ...mapActions('medicalDoctor', ['getAllDoctorsForClinic', 'deleteDoctor', 'updateDoctor']),
    ...mapActions('clinic', ['fetchClinics']),
    initialize() {
      this.fetchClinics();
    },
    findDoctors() {
      this.getAllDoctorsForClinic(this.clinic.id);
    },
    editItem(item) {
      this.editedIndex = this.doctors.indexOf(item);
      this.editedItem = Object.assign(this.editedItem, item);
      this.dialog = true;
    },

    deleteItem(item) {
      this.deleteDoctor(item.id);
    },

    close() {
      this.dialog = false;
      this.$nextTick(() => {
        this.editedItem = Object.assign(this.editedItem, this.defaultItem);
        this.editedIndex = -1;
      });
    },

    save() {
      if (this.editedIndex > -1) {
        const workStart = {
          hour: +this.editedItem.workStart.split(':')[0],
          minute: +this.editedItem.workStart.split(':')[1],
        };
        const workEnd = {
          hour: +this.editedItem.workEnd.split(':')[0],
          minute: +this.editedItem.workEnd.split(':')[1],
        };
        const updatedDoctor = {
          firstName: this.editedItem.firstName,
          lastName: this.editedItem.lastName,
          workStart,
          workEnd,
        };
        this.updateDoctor({ doctorId: this.editedItem.id, updatedDoctor });
      } else {
        this.doctors.push(this.editedItem);
      }
      this.close();
    },
  },
};
</script>
