import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import PrimeVue from 'primevue/config';
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import Dropdown from 'primevue/dropdown';
import FileUpload from 'primevue/fileupload';
import Toolbar from 'primevue/toolbar';
import SplitButton from 'primevue/splitbutton';
import ToastService from 'primevue/toastservice';
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab';
import Toast from 'primevue/toast';
import OrderList from 'primevue/orderlist';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Row from 'primevue/row';                     //optional for row
import Team from '@/components/Team.vue'
import TeamTable from '@/components/Layout/TeamTable.vue'
import Room from '@/components/Layout/Room.vue'


import 'primevue/resources/themes/saga-blue/theme.css'       //theme
import 'primevue/resources/primevue.min.css'                 //core css
import 'primeicons/primeicons.css'                           //icons
import 'primeflex/primeflex.css';
import Card from 'primevue/card';
import Slider from "primevue/slider";
import axios from "axios";

const app = createApp(App);

app.use(createPinia())
app.use(router)
app.use(PrimeVue);
app.use(ToastService);

app.component('Button', Button)
app.component('InputText', InputText)
app.component('FileUpload', FileUpload)
app.component('Dropdown', Dropdown)
app.component('Card', Card)
app.component('Toolbar', Toolbar)
app.component('SplitButton', SplitButton)
app.component('Toast', Toast)
app.component('Accordion', Accordion)
app.component('AccordionTab', AccordionTab)
app.component('OrderList', OrderList)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('Row', Row)
app.component('Team', Team)
app.component('Room', Room)
app.component('TeamTable', TeamTable)
app.component('Slider', Slider)

app.mount('#app')


if (import.meta.env.PROD) {
    const host = window.location.hostname;
    axios.defaults.baseURL = "https://" + host;
}
