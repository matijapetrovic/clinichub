<template>
  <v-container>
    <v-card max-width="1200" max-height="1200" class="mx-auto">
      <v-toolbar
        color="primary"
        dark
        flat
      >
        <v-toolbar-title>Add doctor</v-toolbar-title>
      </v-toolbar>
      <v-card-text>
        <v-form
          ref="form"
        >
          <v-row>
            <v-col
              md="12"
            >
              <v-row>
                <v-col>
                  <v-text-field
                    v-model="firstName"
                    label="First Name"
                    name="firstName"
                    prepend-icon="person"
                    :rules="[requiredRule]"
                    type="text"
                    required
                  ></v-text-field>
                </v-col>
                <v-col>
                  <v-text-field
                    v-model="lastName"
                    label="Last Name"
                    name="lastName"
                    prepend-icon="person"
                    :rules="[requiredRule]"
                    type="text"
                    required
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-select
                    :items="clinics"
                    v-model="clinic"
                    label="Clinic"
                    item-text="name"
                    return-object
                    prepend-icon="mdi-hospital-box-outline"
                  ></v-select>
                </v-col>
                <v-col>
                  <v-select
                    :items="appointmentTypes"
                    v-model="appointmentType"
                    label="Specialization"
                    item-text="name"
                    return-object
                    prepend-icon="mdi-account-cog"
                  ></v-select>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-select
                    :items="dayHours"
                    v-model="from"
                    label="Start of working hours"
                    :rules="[requiredRule]"
                    prepend-icon="mdi-arrow-down-bold-circle-outline"
                  ></v-select>
                </v-col>
                <v-col>
                  <v-select
                    :items="dayHours"
                    v-model="to"
                    label="End of working hours"
                    :rules="[requiredRule]"
                    prepend-icon="mdi-arrow-up-bold-circle-outline"
                  ></v-select>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          @click="submit"
          color="primary"
          name="button"
        >
          Add
        </v-btn>
        <v-spacer></v-spacer>
      </v-card-actions>
      </v-card>
  </v-container>
</template>

<script>
import { mapActions, mapState } from 'vuex';

export default {
  name: 'AddMedicalDoctorForm',
  components: {
  },
  data: () => ({
    email: '',
    appointmentType: '',
    clinic: null,
    password: '',
    confirmPassword: '',
    firstName: '',
    lastName: '',
    address: '',
    city: '',
    country: '',
    telephoneNum: '',
    from: '',
    to: '',
    dayHours: [
      '00:00', '01:00', '02:00', '03:00', '04:00', '05:00',
      '06:00', '07:00', '08:00', '09:00', '10:00', '11:00',
      '12:00', '13:00', '14:00', '15:00', '16:00', '17:00',
      '18:00', '19:00', '20:00', '21:00', '22:00', '23:00',
    ],
  }),
  created() {
    this.initialize();
  },
  methods: {
    ...mapActions('medicalDoctor', ['addMedicalDoctor']),
    ...mapActions('appointmentType', ['fetchAppointmentTypes']),
    ...mapActions('clinic', ['fetchClinics']),
    initialize() {
      this.fetchClinics();
      this.fetchAppointmentTypes();
    },
    submit() {
      if (this.validate()) {
        const workStart = {
          hour: +this.from.split(':')[0],
          minute: +this.from.split(':')[1],
        };
        const workEnd = {
          hour: +this.to.split(':')[0],
          minute: +this.to.split(':')[1],
        };
        this.addMedicalDoctor({
          firstName: this.firstName,
          lastName: this.lastName,
          workStart,
          workEnd,
          specializationId: this.appointmentType.id,
          clinicId: this.clinic.id,
        })
          .then(() => {
            this.clear();
          });
      }
    },
    validate() {
      return this.$refs.form.validate();
    },
    clear() {
      this.$refs.form.reset();
    },
  },
  computed: {
    ...mapState('appointmentType', ['appointmentTypes']),
    ...mapState('clinic', ['clinics']),
    passwordConfirmRule() {
      return () => this.password === this.confirmPassword || 'Passwords must match';
    },
    telephoneNumRule() {
      return (value) => /((\+381)|0)6[0-9]{7,8}/.test(value) || 'Telephone number must be valid';
    },
    requiredRule() {
      return (value) => !!value || 'Required';
    },
  },
};
</script>

<style lang="scss">
p {
  .success{
    color: green;
  }
  .failure {
    color: red;
  }
}
</style>
