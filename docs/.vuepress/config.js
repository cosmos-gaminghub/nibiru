module.exports = {
    title: "Nibiru",
    description: "Documentation for the Cosmos Gaming Hub",
    head: [
        ['link', {rel: 'icon', href: '/assets/logo.png'}],
		],
    markdown: {
        lineNumbers: true,
    },
    plugins: [],
    themeConfig: {
        repo: "cosmos-gaminghub/nibiru",
        editLinks: true,
        docsDir: "docs",
        docsBranch: "master",
        editLinkText: 'Edit this page on Github',
        lastUpdated: true,
        logo: "/assets/logo.png",
        nav: [
            {text: "Website", link: "https://cosmosgaminghub.com", target: "_blank"},
        ],
        sidebarDepth: 2,
        sidebar: [
            {
                title: "Installation",
                collapsable: false,
                children: [
                    ["install/install", "Install"],
                ]
            },
            {
                title: "Testnet",
                collapsable: false,
                children: [
                    ["testnets/fullnode", "Fullnode"],
                    ["testnets/validator", "Validator"],
                ]
            },
            {
                title: "Localnet",
                collapsable: false,
                children: [
                    ["localnets/localnet", "Localnet"],
                    ["localnets/4-node", "4-Node"],
                ]
            },
            {
                title: "Config",
                collapsable: false,
                children: [
                    ["config/service", "Service"],
                ]
            },
        ],
    }
};
