(function() {
    let eles = document.querySelectorAll('img[data-src]')

    if (eles.length === 0) {
        return;
    }

    const parseArray = function(imgs) {
        try {
            return JSON.parse(imgs);
        } catch(err) {
            return [];
        }
    };

    // Load async
    const loadImage = function(url) {
        return new Promise(function(resolv, reject) {
            let img = new Image();
            img.src = url;

            img.onload = function() {
                resolv(this);
            };

            img.onerror = reject;

            img.onabort = reject;
        });
    }

    // Array
    eles = Array.from(eles);

    const runPromisesInSeries = ps => ps.reduce((p, next) => p.then(next), Promise.resolve());

    // Elements
    for (let i = 0, ele; ele = eles[i]; ++i) {
        let promises = []
        // Pictures
        const imgs = parseArray(ele.dataset.src);
        for (let a = 0, img; img = imgs[a]; ++a) {
            promises.push(new Promise(function(resolv) {
                setTimeout(function() {
                    loadImage(img).then(function(imgLoad) {
                        ele.src = imgLoad.src;
                        ele.onload = resolv;

                        ele.classList.add('loaded');
                    }).catch(function() {
                        resolv();
                    });
                }, 300 * a);
            }));
        }

        runPromisesInSeries(promises)
    }
})()