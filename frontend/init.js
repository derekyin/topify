document.addEventListener('DOMContentLoaded', function() {
    try {
      let app = firebase.app();
      let features = ['auth', 'database', 'messaging', 'storage'].filter(feature => typeof app[feature] === 'function');
      console.log('Firebase SDK loaded with' + features.join(', '))
    } catch (e) {
      console.error(e);
    }
  });