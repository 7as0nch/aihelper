import { mountApp } from './mount';

// Auto-mount if running in standalone mode (not imported as library)
// We check if we are in a browser environment and if the default container exists
if (typeof window !== 'undefined' && document.getElementById('app')) {
    mountApp({ routerMode: 'web' });
}
