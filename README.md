# Research Questionnaire

This is a research questionnaire program specifically designed for the research project "Examining the Impact of Cues from Fake Reviews on Consumer Purchase Intention in Online Shopping"

The majority of this project codes are written by GPT, thanks to GPT for his contribution.

## Build from source
To build this from source, you need to have go toolchain version greater than 1.23.2 installed.

To build the program, just simply run the following command:
```bash
make all
```

The program will be compiled and the binary will be placed in the `output` directory.

## Run the program
Before running the program, you need to prepare a `config.json` file in the same directory as the binary. The `config.json` file should have the following structure:

You can refer to the `config.json.example` file for the example of the `config.json` file.

Also, you will need to have MySQL installed on your machine.

Once you have MySQL installed and configured, create a new database for this program:

```sql
CREATE DATABASE research_questionnaire;
```

Then, you can configure the `config.json` file with the database configuration:
```text
"db_con_url": "<db_username>:<db_pass>@tcp(<db_host>:<db_password>)/<db_name>?charset=utf8mb4&parseTime=True&loc=Local",
```

After you get everything set up, use the following command to create database tables:

```bash
./darwin migrate_db
```

And generate some questions for your db:

```bash
./darwin gen_questions
```

Finally, start the server:

```bash
./darwin
```