/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // Diátaxis structure: Tutorials, How-to Guides, Reference, Explanation
  tutorialSidebar: [
    {
      type: 'category',
      label: '📚 Tutorials',
      description: 'Learning-oriented lessons that take you by the hand through a series of steps',
      items: [
        'tutorials/your-first-api-call',
      ],
    },
    {
      type: 'category',
      label: '📖 Reference',
      description: 'Information-oriented technical descriptions of the machinery',
      items: [
        'reference/graphql-schema',
      ],
    },
    {
      type: 'category',
      label: '💡 Explanation',
      description: 'Understanding-oriented discussions that clarify and illuminate topics',
      items: [
        'explanation/microservice-architecture',
      ],
    },
  ],
};

module.exports = sidebars;
