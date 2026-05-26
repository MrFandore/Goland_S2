param(
    [int]$Iterations = 1000,
    [string]$BaseUrl = "http://localhost:8080"
)

$restUrl = "$BaseUrl/v1/tasks"
$graphqlUrl = "$BaseUrl/query"
$graphqlBody = @{
    query = "query { tasks { id } }"
} | ConvertTo-Json -Compress

Write-Host "Speed comparison"
Write-Host "Iterations: $Iterations"
Write-Host "REST URL: $restUrl"
Write-Host "GraphQL URL: $graphqlUrl"

$restMeasure = Measure-Command {
    for ($i = 0; $i -lt $Iterations; $i++) {
        $null = Invoke-WebRequest -UseBasicParsing $restUrl
    }
}

$graphqlMeasure = Measure-Command {
    for ($i = 0; $i -lt $Iterations; $i++) {
        $null = Invoke-WebRequest -UseBasicParsing $graphqlUrl -Method Post -ContentType "application/json" -Body $graphqlBody
    }
}

$restMs = [math]::Round($restMeasure.TotalMilliseconds, 2)
$graphqlMs = [math]::Round($graphqlMeasure.TotalMilliseconds, 2)
$deltaMs = [math]::Round(($restMeasure.TotalMilliseconds - $graphqlMeasure.TotalMilliseconds), 2)

if ($graphqlMeasure.TotalMilliseconds -gt 0) {
    $ratio = [math]::Round(($restMeasure.TotalMilliseconds / $graphqlMeasure.TotalMilliseconds), 2)
} else {
    $ratio = 0
}

[PSCustomObject]@{
    Iterations          = $Iterations
    RestTotalMs         = $restMs
    GraphqlTotalMs      = $graphqlMs
    RestAvgMs           = [math]::Round(($restMeasure.TotalMilliseconds / $Iterations), 4)
    GraphqlAvgMs        = [math]::Round(($graphqlMeasure.TotalMilliseconds / $Iterations), 4)
    DeltaMs             = $deltaMs
    RestDivGraphqlRatio = $ratio
} | Format-List
