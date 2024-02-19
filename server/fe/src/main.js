import { createApp } from 'vue'
import { createPinia } from 'pinia'

import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';
import App from './App.vue'
import router from './router'

import 'primeflex/primeflex.css';
import 'primevue/resources/themes/saga-blue/theme.css'       //theme
import 'primevue/resources/primevue.min.css'                 //core css
import 'primeicons/primeicons.css'                           //icons
import axios from "axios";

const app = createApp(App);

app.use(createPinia())
app.use(router)
app.use(PrimeVue);
app.use(ToastService);

import Button from 'primevue/button'
app.component('Button', Button)

import InputSwitch from 'primevue/inputswitch'
app.component('InputSwitch', InputSwitch)

import SelectButton from 'primevue/selectbutton'
app.component('SelectButton', SelectButton)

import InputNumber from 'primevue/inputnumber'
app.component('InputNumber', InputNumber)

import RadioButton from 'primevue/radiobutton'
app.component('RadioButton', RadioButton)

import InputText from 'primevue/inputtext'
app.component('InputText', InputText)

import FileUpload from 'primevue/fileupload'
app.component('FileUpload', FileUpload)

import Dropdown from 'primevue/dropdown'
app.component('Dropdown', Dropdown)

import Card from 'primevue/card'
app.component('Card', Card)

import Toolbar from 'primevue/toolbar'
app.component('Toolbar', Toolbar)

import Badge from 'primevue/badge'
app.component('Badge', Badge)

import SplitButton from 'primevue/splitbutton'
app.component('SplitButton', SplitButton)

import Toast from 'primevue/toast'
app.component('Toast', Toast)

import Accordion from 'primevue/accordion'
app.component('Accordion', Accordion)

import TabMenu from 'primevue/tabmenu'
app.component('TabMenu', TabMenu)

import Divider from 'primevue/divider'
app.component('Divider', Divider)

import Panel from 'primevue/panel'
app.component('Panel', Panel)

import Splitter from 'primevue/splitter'
app.component('Splitter', Splitter)

import SplitterPanel from 'primevue/splitterpanel'
app.component('SplitterPanel', SplitterPanel)

import ToggleButton from 'primevue/togglebutton'
app.component('ToggleButton', ToggleButton)

import AccordionTab from 'primevue/accordiontab'
app.component('AccordionTab', AccordionTab)

import OrderList from 'primevue/orderlist'
app.component('OrderList', OrderList)

import DataTable from 'primevue/datatable'
app.component('DataTable', DataTable)

import Column from 'primevue/column'
app.component('Column', Column)

import Row from 'primevue/row'
app.component('Row', Row)

import '@vueform/slider/themes/default.css'
import Slider from '@vueform/slider'
app.component('Slider', Slider)

import Checkbox from 'primevue/checkbox'
app.component('Checkbox', Checkbox)

// Local components
import Team from '@/components/Team.vue'
app.component('Team', Team)

import Room from '@/components/Layout/Room.vue'
app.component('Room', Room)

import TeamTable from '@/components/Layout/TeamTable.vue'
app.component('TeamTable', TeamTable)

app.mount('#app')

if (import.meta.env.PROD) {
    const host = window.location.hostname;
    axios.defaults.baseURL = "https://" + host;
}
