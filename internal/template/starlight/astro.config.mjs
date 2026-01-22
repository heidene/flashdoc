// Learn more: https://starlight.astro.build/
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  integrations: [
    starlight({
      title: '{{SITE_TITLE}}',
      defaultLocale: 'en',
    }),
  ],
});
