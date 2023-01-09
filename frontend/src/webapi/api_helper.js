export default {
    createUrl(...args) {
        console.log(args);
        let url = process.env.BACKEND_URL;
        console.log(url);
        args.forEach((arg) => {url = `${url}/${arg}`});
        console.log(url);
        return url;
    },
    /**
     * @param {string} url 
     * @param {function(object)} onResponse 
     * @param {function(object)} onError 
     */
    get(url, onResponse, onError) {
        fetch(url)
            .then(response => {
                response.json().then(data => {
                    if (response.ok) {
                        onResponse(data);
                        return;
                    }
                    onError(Error(response.statusText));
                })

            })
            .catch(onError)
    },
    /**
     * @param {string} url
     * @param {object} data
     * @param {function(object)} onResponse
     * @param {function(object)} onError
     */
    post(url, data, onResponse, onError) {
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            response.json().then(data => {
                if (response.ok) {
                    onResponse(data);
                    return;
                }
                onError(Error(response.statusText));
            })

        })
        .catch(onError)
    }
}