<template>
  <v-container>
    <v-row
      class="mb-6"
      no-gutters
    >
    <v-col
      cols="12"
      lg="12"
    >
    <v-data-table
      :headers="headers"
      :items="clinics"
      :search="search"
      class="elevation-1"
    >
      <template v-slot:top>
        <v-toolbar flat color="white">
          <v-toolbar-title>Clinics</v-toolbar-title>
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
                <v-container>
                  <v-row>
                    <v-col cols="12" sm="6" md="6">
                      <v-text-field v-model="editedItem.name" label="Name">
                      </v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="6">
                      <v-text-field v-model="editedItem.address.addressLine" label="Address">
                      </v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="6">
                      <v-text-field v-model="editedItem.address.city" label="City">
                      </v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="6">
                       <CountrySelect
                          v-model="editedItem.address.country"
                        />
                    </v-col>
                    <v-col cols="12" sm="12" md="12">
                      <v-textarea v-model="editedItem.description" label="Description">
                      </v-textarea>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
                <v-btn color="blue darken-1" text @click="save">Save</v-btn>
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
        <v-btn color="primary" @click="initialize">Reset</v-btn>
      </template>
    </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import CountrySelect from '@/app/country/_components/CountrySelect.vue';

export default {
  components: {
    CountrySelect,
  },
  data: () => ({
    dialog: false,
    search: '',
    headers: [
      {
        text: 'Name',
        align: 'start',
        value: 'name',
      },
      { text: 'Address', value: 'address.addressLine' },
      { text: 'City', value: 'address.city' },
      { text: 'Country', value: 'address.country' },
      { text: 'Actions', value: 'actions', sortable: false },
    ],
    editedIndex: -1,
    editedItem: {
      name: '',
      description: '',
      address: {
        addressLine: '',
        city: '',
        country: '',
      },
    },
    defaultItem: {
      name: '',
      description: '',
      address: {
        addressLine: '',
        city: '',
        country: '',
      },
    },
  }),

  computed: {
    ...mapState('clinic', ['clinics']),
    formTitle() {
      return this.editedIndex === -1 ? 'New Item' : 'Edit Item';
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
    ...mapActions('clinic', ['fetchClinics', 'updateClinic']),
    initialize() {
      this.fetchClinics();
    },
    editItem(item) {
      this.editedIndex = this.clinics.indexOf(item);
      this.editedItem = JSON.parse(JSON.stringify(item));
      this.dialog = true;
    },

    deleteItem(item) {
      this.deleteClinic(item.id);
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
        const updatedClinic = {
          name: this.editedItem.name,
          description: this.editedItem.description,
          address: {
            addressLine: this.editedItem.address.addressLine,
            city: this.editedItem.address.city,
            country: this.editedItem.address.country,
          },
        };
        this.updateClinic({ clinicId: this.editedItem.id, updatedClinic });
      } else {
        this.clinics.push(this.editedItem);
      }
      this.close();
    },
  },
};
</script>
