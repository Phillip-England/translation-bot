# translation-bot

Does your team use groupme? Are you constantly having to use google translate to ensure everyone on your team can read messages? Introducing: translation-bot! A quick and easy way to auto-translate groupme messages using deepl.com's translation API.

## Usage

translation-bot enables your groupme memeber to auto-translate their messages by prefixing a message with $spanish or $english depending on what lanugage they want to translate their message to. Codebase can be modified to fit more languages if needed.

## Installation

1. Clone the repository

```bash
git clone https://github.com/phillip-england/translation-bot <target-directory>
```

2. Install Go Packages

```bash
go mod tidy
```

3. Create a .env in the root of the project

```bash
touch .env
```

4. Add the required environment varialbes to .env

```bash
DEEPL_API_KEY=your deepl.com api key for translation (they offer a free tier)
PORT=port to serve on
YOUR_BOT_ID=your groupme bot ID (more on this in next section)
YOUR_GROUP_ID=your groupme group ID (more on this in next section)
```

5. Modify the following section in main.go to fit your requirements. In my case, I am using this bot on multiple different groupme channels.

```go
var groupmeBotID string
if groupID == os.Getenv("TESTING_GROUP_ID") {
    groupmeBotID = os.Getenv("TESTING_BOT_ID")
}
if groupID == os.Getenv("SOUTHROADS_LEADERSHIP_GROUP_ID") {
    groupmeBotID = os.Getenv("SOUTHROADS_LEADERSHIP_BOT_ID")
}
if groupID == os.Getenv("KITCHEN_LEADERSHIP_GROUP_ID") {
    groupmeBotID = os.Getenv("KITCHEN_LEADERSHIP_BOT_ID")
}
if groupID == os.Getenv("FOH_OPERATIONS_GROUP_ID") {
    groupmeBotID = os.Getenv("FOH_OPERATIONS_BOT_ID")
}
if groupID == os.Getenv("SUPPLY_ORDER_GROUP_ID") {
    groupmeBotID = os.Getenv("SUPPLY_ORDER_BOT_ID")
}
```

To fit based on the .env file I provided above, modify to the following:

```go
var groupmeBotID string
if groupID == os.Getenv("YOUR_GROUP_ID") {
    groupmeBotID = os.Getenv("YOUR_BOT_ID")
}
```

Simply add more .env varialbes to add additional groups your bot can post to.

6. Serve the Application

```bash
go run main.go
```

## Third-Party Documentation

For more information on how to use the groupme api or the deepl translation api, check the following links:

- [Groupme API](https://dev.groupme.com/)
- [Deepl API](https://www.deepl.com/translator) 

