# `ledger-grievance` opentofu infrastructure

declarative infrastructure as code for deploying and observing `ledger-grievance` with radical kindness and state reconciliation

## the four pillars of declarative grievance healing
1. **core namespace and secret generation** via `main.tf` creating the `grievance-ledger` namespace and generating dynamic `random_password` credentials because your security and emotional safety deserve entropy
2. **mariadb persistence layer** via `database.tf` provisioning the `mariadb-service` and `mariadb` statefulset and automatically mounting `../db/init.sql` into `/docker-entrypoint-initdb.d/init.sql` via a configmap so our database schema remains unified and whole
3. **scalable go api and horizontal pod autoscaling** via `api.tf` provisioning the `ledger-api` deployment `ledger-service` clusterip service and `ledger-hpa` horizontal pod autoscaler with `/health` and `/health/deep` probes so that kindness always scales proportionally with institutional demand
4. **wholesome observability and grafana dashboards** via `observability.tf` and `dashboards/wholesome.json` deploying prometheus to scrape `ledger-service` alongside grafana preloaded with our wholesome grievance telemetry dashboard tracking `wholesome_notes_total` and `gopher_wisdom_total` and active goroutines holding space for your feelings

## howto
### prerequisites
- opentofu 1.6+ or terraform 1.6+
- a kubernetes cluster ready to receive your grievances

### usage via makefile
from the repository root you can invoke our auspicious tool agnostic makefile targets which support both `tofu` and `terraform` via `IAC` and `IAC_DIR`
```bash
make iac-init
make iac-fmt
make iac-validate
make iac-plan
```

### direct opentofu usage
1. initialize the provider plugins without fear
   ```bash
   cd tofu
   tofu init
   ```

2. validate that your declarative intentions are pure and syntactically sound
   ```bash
   tofu validate
   ```

3. preview the reconciliation of your infrastructure state
   ```bash
   tofu plan -var="environment=staging"
   ```

4. apply the configuration to manifest wholesome observability directly into a running local grafana instance
   ```bash
   tofu apply -var="enable_grafana_provider_resources=true"
   ```
