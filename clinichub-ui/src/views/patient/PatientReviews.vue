<template>
  <div>
    <h1>Ratings</h1>
    <v-row>
      <v-col>
        <h2>Clinics</h2>
        <Reviews
          :items="clinicsForReview"
          label="Select a clinic"
          @rated="rateClinic($event)"
        />
      </v-col>
      <v-col>
        <h2>Doctors</h2>
        <Reviews
          :items="doctorsForReview"
          label="Select a doctor"
          @rated="rateDoctor($event)"
        />
      </v-col>
    </v-row>
  </div>
</template>

<script>
import Reviews from '@/app/reviews/Reviews.vue';
import { mapActions, mapState } from 'vuex';

export default {
  name: 'PatientReviews',
  components: {
    Reviews,
  },
  computed: {
    ...mapState('reviews', ['clinicsForReview', 'doctorsForReview']),
  },
  methods: {
    ...mapActions('reviews', ['fetchClinicsForReview', 'fetchDoctorsForReview', 'addClinicReview', 'addDoctorReview']),
    rateClinic(rating) {
      this.addClinicReview({ clinicId: rating.id, rating: rating.rating });
    },
    rateDoctor(rating) {
      this.addDoctorReview({ doctorId: rating.id, rating: rating.rating });
    },
  },
  mounted() {
    this.fetchClinicsForReview();
    this.fetchDoctorsForReview();
  },
};
</script>
