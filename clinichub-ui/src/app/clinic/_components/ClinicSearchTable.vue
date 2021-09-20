<template>
  <v-card>
    <v-card-title>
      Clinics
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Filter"
        single-line
        hide-details
      >

      </v-text-field>
    </v-card-title>
    <v-data-table
    :headers="headers"
    :items="items"
    :items-per-page="5"
    :search="search"
    class="elevation-1"
  >
    <template v-slot:body="{ items }">
        <tbody>
          <tr v-for="item in items" :key="item.name">
            <td>
              <v-btn
                color="red lighten-2"
                outlined
                rounded
                @click.stop="openClinicProfile(item)"
              >
                {{ item.name }}
              </v-btn>
            </td>
            <td>{{ item.address.addressLine }}</td>
            <td>{{ item.address.city }}</td>
            <td>{{ item.address.country }}</td>
            <td>{{ parseFloat(item.rating.rating).toFixed(2) }} ({{ item.rating.count}})</td>
            <td>{{ item.price ? item.price + ' €' : 'N/A' }}</td>
          </tr>
        </tbody>
      </template>
  </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: 'ClinicSearchTable',
  data: () => ({
    headers: [
      {
        text: 'Name',
        align: 'start',
        sortable: true,
        value: 'name',
      },
      { text: 'Address', value: 'address' },
      { text: 'City', value: 'city' },
      { text: 'Country', value: 'country' },
      { text: 'Rating', value: 'rating' },
      { text: 'Price', value: 'appointmentPrice' },
    ],
    search: '',
  }),
  props: {
    items: {
      type: Array,
      required: true,
    },
  },
  methods: {
    openClinicProfile(clinic) {
      this.$router.push(`/clinic/${clinic.id}`);
    },
  },
};
</script>
