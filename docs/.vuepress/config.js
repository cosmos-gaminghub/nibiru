module.exports = {
    title: "Nibiru",
    description: "Documentation for the Cosmos Gaming Hub",
    head: [],
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
                title: "Testnets",
                collapsable: false,
                children: [
                    ["testnets/overview", "Overview"]
                ]
            },
            {
                title: "Fullnode",
                collapsable: false,
                children: [
                    ["fullnode/install", "Install"]
                ]
            },
            {
                title: "Validators",
                collapsable: false,
                children: [
                    ["validators/setup", "Setup"]
                ]
            }
        ],
    }
};
