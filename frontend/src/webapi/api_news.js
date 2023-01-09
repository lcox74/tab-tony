import api from './api_helper';

const BACKEND_URL = "http://localhost:3000/"

export default {

    validateAccess(accesskey, onSuccess, onError) {
        api.get(BACKEND_URL + "news/" + accesskey, onSuccess, onError);
    },

    createNewsPost(post, onSuccess, onError) {
        api.post(BACKEND_URL + "news", post, onSuccess, onError);
    }
}
