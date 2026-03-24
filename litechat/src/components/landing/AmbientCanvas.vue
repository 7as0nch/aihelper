<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';

const props = defineProps<{
  theme: 'light' | 'dark';
  progress: number;
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
const isEnabled = ref(false);

let stop = () => {};

onMounted(async () => {
  if (typeof window === 'undefined') {
    return;
  }

  const media = window.matchMedia('(prefers-reduced-motion: reduce)');
  if (media.matches || !canvasRef.value) {
    return;
  }

  try {
    const THREE = await import('three');
    const canvas = canvasRef.value;
    if (!canvas) {
      return;
    }

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(42, 1, 0.1, 40);
    camera.position.z = 12;

    const renderer = new THREE.WebGLRenderer({
      canvas,
      antialias: true,
      alpha: true,
      powerPreference: 'high-performance',
    });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 1.5));

    const group = new THREE.Group();
    scene.add(group);

    const particleCount = 180;
    const positions = new Float32Array(particleCount * 3);
    for (let index = 0; index < particleCount; index += 1) {
      const offset = index * 3;
      positions[offset] = (Math.random() - 0.5) * 18;
      positions[offset + 1] = (Math.random() - 0.5) * 16;
      positions[offset + 2] = (Math.random() - 0.5) * 8;
    }

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));

    const particles = new THREE.Points(
      geometry,
      new THREE.PointsMaterial({
        color: props.theme === 'dark' ? '#7dd3fc' : '#60a5fa',
        size: 0.1,
        transparent: true,
        opacity: 0.72,
        depthWrite: false,
      }),
    );
    group.add(particles);

    const halo = new THREE.Mesh(
      new THREE.TorusGeometry(3.5, 0.08, 16, 100),
      new THREE.MeshBasicMaterial({
        color: props.theme === 'dark' ? '#60a5fa' : '#38bdf8',
        transparent: true,
        opacity: 0.28,
      }),
    );
    halo.rotation.x = 1.25;
    group.add(halo);

    const accent = new THREE.Mesh(
      new THREE.TorusGeometry(5.1, 0.05, 12, 90),
      new THREE.MeshBasicMaterial({
        color: props.theme === 'dark' ? '#c4b5fd' : '#6366f1',
        transparent: true,
        opacity: 0.18,
      }),
    );
    accent.rotation.x = 1.06;
    accent.rotation.z = -0.35;
    group.add(accent);

    const updateTheme = (nextTheme: 'light' | 'dark') => {
      const particleMaterial = particles.material as typeof particles.material & { color: { set: (color: string) => void } };
      const haloMaterial = halo.material as typeof halo.material & { color: { set: (color: string) => void } };
      const accentMaterial = accent.material as typeof accent.material & { color: { set: (color: string) => void } };

      particleMaterial.color.set(nextTheme === 'dark' ? '#7dd3fc' : '#60a5fa');
      haloMaterial.color.set(nextTheme === 'dark' ? '#60a5fa' : '#38bdf8');
      accentMaterial.color.set(nextTheme === 'dark' ? '#c4b5fd' : '#6366f1');
    };

    const resize = () => {
      if (!canvas.parentElement) {
        return;
      }
      const { clientWidth, clientHeight } = canvas.parentElement;
      if (!clientWidth || !clientHeight) {
        return;
      }
      camera.aspect = clientWidth / clientHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(clientWidth, clientHeight, false);
    };

    resize();
    window.addEventListener('resize', resize);

    const themeStop = watch(
      () => props.theme,
      (nextTheme) => {
        updateTheme(nextTheme);
      },
      { immediate: true },
    );

    let frameId = 0;
    const clock = new THREE.Clock();
    const tick = () => {
      const elapsed = clock.getElapsedTime();
      group.rotation.z = elapsed * 0.045 + props.progress * 0.75;
      group.rotation.x = Math.sin(elapsed * 0.22) * 0.08;
      group.position.y = (0.5 - props.progress) * 1.4;
      particles.rotation.y = elapsed * 0.03;
      halo.rotation.z = elapsed * 0.06;
      accent.rotation.y = elapsed * 0.05;
      renderer.render(scene, camera);
      frameId = window.requestAnimationFrame(tick);
    };

    isEnabled.value = true;
    tick();

    stop = () => {
      themeStop();
      window.cancelAnimationFrame(frameId);
      window.removeEventListener('resize', resize);
      geometry.dispose();
      particles.material.dispose();
      halo.geometry.dispose();
      halo.material.dispose();
      accent.geometry.dispose();
      accent.material.dispose();
      renderer.dispose();
    };
  } catch (_error) {
    isEnabled.value = false;
  }
});

onBeforeUnmount(() => {
  stop();
});
</script>

<template>
  <div class="ambient-wrap" aria-hidden="true">
    <canvas ref="canvasRef" class="ambient-canvas" :class="{ 'is-visible': isEnabled }" />
    <div class="ambient-fallback" />
  </div>
</template>

<style scoped>
.ambient-wrap {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}

.ambient-canvas,
.ambient-fallback {
  position: absolute;
  inset: 0;
}

.ambient-canvas {
  opacity: 0;
  transition: opacity 0.5s ease;
}

.ambient-canvas.is-visible {
  opacity: 0.82;
}

.ambient-fallback {
  background:
    radial-gradient(circle at 16% 18%, rgba(96, 165, 250, 0.26), transparent 28%),
    radial-gradient(circle at 82% 16%, rgba(191, 219, 254, 0.42), transparent 24%),
    radial-gradient(circle at 52% 78%, rgba(165, 180, 252, 0.24), transparent 28%);
  filter: blur(12px);
  opacity: 0.75;
}

:global(.theme-dark) .ambient-fallback {
  background:
    radial-gradient(circle at 18% 16%, rgba(29, 78, 216, 0.28), transparent 30%),
    radial-gradient(circle at 80% 18%, rgba(56, 189, 248, 0.2), transparent 24%),
    radial-gradient(circle at 58% 76%, rgba(196, 181, 253, 0.12), transparent 28%);
}
</style>