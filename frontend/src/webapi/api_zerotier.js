import api from './api_helper';

const BACKEND_URL = "http://light-candle.bnr.la:3000/"
// const BACKEND_URL = "http://localhost:3000/"

export default {

    validateAccess(accesskey, onSuccess, onError) {
        api.get(BACKEND_URL + "zerotier/" + accesskey, onSuccess, onError);
    },

    getNetworkMembers(accesskey, onSuccess, onError) {
        api.get(BACKEND_URL + "zerotier/" + accesskey + "/network", onSuccess, onError);
    }
}
