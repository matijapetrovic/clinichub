<template>
  <v-container>
    <v-row>
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
        color="primary"
        class="mt-5 ml-5"
        @click="findPrices">
        Find Prices
      </v-btn>
    </v-col>
    </v-row>
    <v-data-table
      :headers="headers"
      :items="prices"
      class="elevation-1"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>Appointment Prices</v-toolbar-title>
        <v-btn :disabled="!enableAdd"
          @click="newPrice"
          class="ml-5" color="secondary">Add Price</v-btn>
        </v-toolbar>
      </template>
      <template v-slot:item.price="{ item }">
        <div v-if="item.price === undefined" class="my-2">
          <v-btn
          rounded
          small
          color="success"
          @click="setItem(item)"
          >
            Add price
          </v-btn>
        </div>
        <div v-else>
          <v-btn
          rounded
          small
          color="success"
          @click="setItem(item)"
          >
            Edit price ({{item.price}} $)
          </v-btn>
        </div>
      </template>
      <template v-slot:no-data>
        No prices available
      </template>
    </v-data-table>
    <v-dialog
      v-model="dialog"
      max-width="420"
    >
      <v-card>
        <v-card-title class="headline">
          Appointment type price
        </v-card-title>
        <v-card-text>
          <v-text-field
          v-if="editedItem.appointmentType && !isNew"
          v-model="editedItem.appointmentType.name"
          label="Appointment type"
          readonly>
          </v-text-field>
          <v-select
            v-if="isNew"
            :items="noPriceAppointmentTypes"
            v-model="editedItem.appointmentType"
            label="Appointment Type"
            item-text="name"
            return-object
          ></v-select>
          <v-text-field
            v-model="editedItem.price"
            label="Price"
            type="number"
            min="1"
            :rules="[requiredRule,]"
          >
          </v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="red darken-1"
            text
            @click="closeDialog()"
          >
            Close
          </v-btn>

          <v-btn
            color="green darken-1"
            text
            @click="addAppointmentTypePricePrice()"
          >
            Add
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import { mapActions, mapState } from 'vuex';

export default {
  data: () => ({
    dialog: false,
    headers: [
      {
        text: 'Appointment type',
        align: 'start',
        sortable: false,
        value: 'appointmentType.name',
      },
      { text: 'Price', value: 'price' },
    ],
    items: [],
    editedIndex: -1,
    editedItem: {
      appointmentType: null,
      price: null,
    },
    defaultItem: {
      appointmentType: null,
      price: null,
    },
    editedItemFirstValue: null,
    price: null,
    clinic: null,
    isNew: false,
    enableAdd: false,
  }),

  mounted() {
    this.fetchClinics();
    this.fetchAppointmentTypes();
  },
  methods: {
    ...mapActions('clinic', ['fetchPrices', 'addPrice', 'updatePrice', 'fetchClinics']),
    ...mapActions('appointmentType', ['fetchAppointmentTypes']),
    findPrices() {
      this.fetchPrices(this.clinic.id);
      this.enableAdd = true;
    },
    newPrice() {
      this.isNew = true;
      this.dialog = true;
    },
    editItem(item) {
      this.isNew = false;
      this.editedIndex = this.items.map((e) => e.id).indexOf(item.id);
      this.price = item;
      this.dialog = true;
    },
    deleteItem(item) {
      const index = this.items.indexOf(item);
      window.confirm('Are you sure you want to delete this item?');
      this.items.splice(index, 1);
    },
    setItem(item) {
      this.isNew = false;
      this.editedItem = { ...item };
      this.price = item.price;
      this.dialog = true;
    },
    setItems() {
      this.items.length = 0;
      this.appointmentTypes.forEach((item) => this.items.push({
        appointmentType: item,
        price: this.prices[(item.id).toString()],
      }));
    },
    addAppointmentTypePricePrice() {
      if (this.isNew) {
        this.addPrice({
          clinicId: this.clinic.id,
          price: {
            appointmentTypeId: this.editedItem.appointmentType.id,
            price: +this.editedItem.price,
          },
        })
          .then(() => {
            this.fetchPrices(this.clinic.id);
          });
        this.closeDialog();
      } else {
        this.updatePrice({
          clinicId: this.clinic.id,
          price: {
            appointmentTypeId: this.editedItem.appointmentType.id,
            price: +this.editedItem.price,
          },
        })
          .then(() => {
            this.fetchPrices(this.clinic.id);
          });
        this.closeDialog();
      }
    },
    closeDialog() {
      this.editedItem = Object.assign(this.editedItem, this.defaultItem);
      this.dialog = false;
    },
  },
  computed: {
    ...mapState('clinic', ['prices', 'clinics']),
    ...mapState('appointmentType', ['appointmentTypes']),
    requiredRule() {
      return (value) => !!value || 'Required';
    },
    noPriceAppointmentTypes() {
      console.log(this.prices);
      return this.appointmentTypes.filter((appType) => !this.prices.some(
        (price) => price.appointmentType.id === appType.id,
      ));
    },
  },
};
</script>
