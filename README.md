by [andrewsaputra](https://github.com/andrewsaputra)

# Sample Rest App With Golang

I made this project to learn Go language and experiment with its features. Also to experiment with designing a REST compliant APIs.

Experimentations covered :
- Web Frameworks, e.g. : GIN
- REST compliant APIs
- Data store CRUD operations with :
    - Ephemeral In Memory
    - MongoDB
    - DynamoDB
- CI / CD integrations with AWS Services : [Terraform](https://github.com/andrewsaputra/aws-sandbox)

### Tech Stacks
- [Go 1.21.4](https://go.dev/doc/install)
- [Gin Web Framework](https://gin-gonic.com/)
- [MongoDB 6.0.11](https://www.mongodb.com/docs/v6.0/tutorial/install-mongodb-on-ubuntu/)
- [DynamoDB 2.1.0](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html)


### API Endpoints

<table>
    <thead>
        <tr>
            <th>Path</th>
            <th>Method</th>
            <th>Payload Example</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>/status</code></td>
            <td>GET</td>
            <td></td>
            <td>Application Status / Health Check</td>
        </tr>
        <tr>
            <td><code>/albums</code></td>
            <td>GET</td>
            <td></td>
            <td>Retrieve all albums records</td>
        </tr>
        <tr>
            <td><code>/albums</code></td>
            <td>POST</td>
            <td>
                <details>
                    <summary>example</summary>
                    <pre><code>{"title":"song 1", "artist":"singer A", "price":9.99}</code></pre>
                </details>
            </td>
            <td>Insert new album record</td>
        </tr>
        <tr>
            <td><code>/albums/{id}</code></td>
            <td>GET</td>
            <td></td>
            <td>Retrieve specific album record</td>
        </tr>
        <tr>
            <td><code>/albums/{id}</code></td>
            <td>PUT</td>
            <td>
                <details>
                    <summary>example</summary>
                    <pre><code>{"title":"song 1", "artist":"singer A", "price":9.99}</code></pre>
                </details>
            </td>
            <td>Replace album record's properties</td>
        </tr>
        <tr>
            <td><code>/albums/{id}</code></td>
            <td>PATCH</td>
            <td>
                <details>
                    <summary>example</summary>
                    <pre><code>{"price":11.11}</code></pre>
                </details>
            </td>
            <td>Partial update to album record's properties</td>
        </tr>
        <tr>
            <td><code>/albums/{id}</code></td>
            <td>DELETE</td>
            <td></td>
            <td>Delete album record</td>
        </tr>
    </tbody>
</table>

